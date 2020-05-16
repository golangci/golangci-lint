import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'gatsby';
import { MdKeyboardArrowRight, MdKeyboardArrowLeft } from 'react-icons/md';

import { Container, Post } from './styles';

export default function Docs({ prev, next }) {
  return (
    <Container>
      {prev && (
        <Post isLeft>
          <Link to={prev.link}>
            <MdKeyboardArrowLeft />
            <div>
              <p>Prev</p>
              <h3>{prev.label}</h3>
            </div>
          </Link>
        </Post>
      )}
      {next && (
        <Post>
          <Link to={next.link}>
            <div>
              <p>Next</p>
              <h3>{next.label}</h3>
            </div>
            <MdKeyboardArrowRight />
          </Link>
        </Post>
      )}
    </Container>
  );
}

Docs.propTypes = {
  prev: PropTypes.shape({
    label: PropTypes.string,
    link: PropTypes.string,
  }),
  next: PropTypes.shape({
    label: PropTypes.string,
    link: PropTypes.string,
  }),
};

Docs.defaultProps = {
  prev: null,
  next: null,
};
