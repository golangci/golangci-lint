import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'gatsby';

export default function InternalLink({ link, label }) {
  return (
    <Link to={link} activeClassName="active-link">
      {label}
    </Link>
  );
}

InternalLink.propTypes = {
  link: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
};
