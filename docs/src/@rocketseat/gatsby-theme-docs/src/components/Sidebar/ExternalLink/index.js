import React from 'react';
import PropTypes from 'prop-types';
import { FiExternalLink } from 'react-icons/fi';

export default function ExternalLink({ link, label }) {
  return (
    <a href={link} rel="noopener noreferrer">
      {label}
      <FiExternalLink
        style={{ width: '16px', height: '16px', marginLeft: '10px' }}
      />
    </a>
  );
}

ExternalLink.propTypes = {
  link: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
};
