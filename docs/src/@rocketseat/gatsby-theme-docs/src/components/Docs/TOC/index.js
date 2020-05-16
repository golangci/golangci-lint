import React from 'react';
import PropTypes from 'prop-types';

import slug from '../../../util/slug';

import { Container } from './styles';

export default function TableOfContents({ headings }) {
  if (headings && headings.length !== 0) {
    return (
      <Container>
        <h2>Table of Contents</h2>
        <nav>
          <ul>
            {headings
              .filter(heading => heading.depth === 2)
              .map(heading => (
                <li key={heading.value}>
                  <a href={`#${slug(heading.value)}`}>{heading.value}</a>
                </li>
              ))}
          </ul>
        </nav>
      </Container>
    );
  }

  return null;
}

TableOfContents.propTypes = {
  headings: PropTypes.array,
};

TableOfContents.defaultProps = {
  headings: null,
};
