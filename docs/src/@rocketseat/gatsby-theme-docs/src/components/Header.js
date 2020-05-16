import React from 'react';
import PropTypes from 'prop-types';
import Headroom from 'react-headroom';
import styled from '@emotion/styled';
import { GiHamburgerMenu } from 'react-icons/gi';
import { useStaticQuery, graphql } from 'gatsby';

const Container = styled.header`
  display: flex;
  justify-content: flex-start;
  align-items: center;

  height: 60px;
  padding: 0 24px;
  background: #fff;

  transition: transform 0.5s;

  transform: translate3d(
    ${({ isMenuOpen }) => (isMenuOpen ? '240px' : '0')},
    0,
    0
  );

  h2 {
    margin: 0;
    border: none;
    padding: 0;
    font-size: 18px;
    color: #000;
  }

  button {
    border: none;
    background: #fff;
    cursor: pointer;
    margin-right: 16px;
  }

  @media (min-width: 780px) {
    display: none;
  }
`;

export default function Header({ handleMenuOpen, isMenuOpen }) {
  const { site } = useStaticQuery(
    graphql`
      query {
        site {
          siteMetadata {
            siteTitleShort
          }
        }
      }
    `,
  );

  const { siteTitleShort } = site.siteMetadata;

  return (
    <Headroom>
      <Container isMenuOpen={isMenuOpen}>
        <button
          aria-label="Open sidebar"
          type="button"
          onClick={handleMenuOpen}
        >
          <GiHamburgerMenu size={23} />
        </button>
        <h2>{siteTitleShort}</h2>
      </Container>
    </Headroom>
  );
}

Header.propTypes = {
  handleMenuOpen: PropTypes.func.isRequired,
  isMenuOpen: PropTypes.bool.isRequired,
};
