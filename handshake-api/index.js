const url = require("url");
const express = require("express");
const proxy = require("express-http-proxy");
const { NodeClient } = require("hs-client");

const host = process.env.HOSTNAME || "0.0.0.0";
const port = Number(process.env.PORT) || 3100;

const hsdNetworkType = process.env.HSD_NETWORK || "regtest";
const hsdHost = process.env.HSD_HOST || "localhost";
const hsdPort = Number(process.env.HSD_PORT) || 12037;
const hsdApiKey = process.env.HSD_API_KEY || "foo";

const clientOptions = {
  network: hsdNetworkType,
  host: hsdHost,
  port: hsdPort,
  apiKey: hsdApiKey,
};
const client = new NodeClient(clientOptions);

const startsWithSkylinkRegExp = /^[a-zA-Z0-9_-]{46}/;

const getDomainRecords = async (name) => {
  const response = await client.execute("getnameresource", [name]);
  const records = response?.records ?? null;

  console.log(`${name} => ${JSON.stringify(records)}`);

  return records;
};

const findSkylinkRecord = (records) => {
  return records?.find(({ txt }) => txt?.some((entry) => isValidSkylink(entry)));
};

const getSkylinkFromRecord = (record) => {
  return record?.txt?.find((entry) => isValidSkylink(entry));
};

const resolveDomainHandler = async (req, res) => {
  try {
    const records = await getDomainRecords(req.params.name);
    if (!records) return res.status(404).send(`No records found for ${req.params.name}`);

    const record = findSkylinkRecord(records);
    if (!record) throw new Error(`No skylink found in dns records of ${req.params.name}`);

    const skylink = getSkylinkFromRecord(record);
    return res.json({ skylink });
  } catch (error) {
    res.status(500).send(`Handshake error: ${error.message}`);
  }
};

// Checks if the given string is a valid Sia link.
function isValidSkylink(link) {
  if (!link || link.length === 0) {
    return false;
  }
  return Boolean(link.match(startsWithSkylinkRegExp));
}

const server = express();

server.use(
  "/hns/:name",
  proxy("localhost:8080", {
    // eslint-disable-next-line no-unused-vars
    userResHeaderDecorator(headers, userReq, userRes, proxyReq, proxyRes) {
      if (headers.location && headers.location.match(startsWithSkylinkRegExp)) {
        headers.location = headers.location.replace(startsWithSkylinkRegExp, `/hns/${userReq.params.name}`);
      }

      return headers;
    },
    proxyReqPathResolver: async (req) => {
      const records = await getDomainRecords(req.params.name);
      if (!records) throw new Error(`No records found for ${req.params.name}`);

      const record = findSkylinkRecord(records);
      if (!record) throw new Error(`No skylink found in dns records of ${req.params.name}`);

      const skylink = getSkylinkFromRecord(record);
      const basepath = url.resolve("/", skylink); // make the url absolute
      const subpath = req.url.slice(1); // drop the leading slash

      // if the skylink from handshake does not contain a subpath but subpath
      // is defined in request, join the skylink and subpath together (do not
      // use url.resolve because it will replace skylink with subapth thinking
      // it is relative)
      if (skylink.length === 46 && subpath) {
        return `${basepath}/${subpath}`;
      }

      return url.resolve(basepath, subpath);
    },
  })
);

server.get("/hnsres/:name", resolveDomainHandler);

server.listen(port, host, (error) => {
  if (error) throw error;

  console.info(`API will look for HSD Server at ${hsdHost}:${hsdPort}, "${hsdNetworkType}" network.`);
  console.info(`Server listening at http://${host}:${port}`);
});