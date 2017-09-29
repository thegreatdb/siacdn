import Router from 'next/router';

const clickRouter = path => {
  return ev => {
    Router.push(path);
    ev.preventDefault();
    ev.stopPropagation();
  };
};

export default clickRouter;
