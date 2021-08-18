import styled from '@emotion/styled';
import { css } from '@emotion/react';

export const Main = styled.main`
  padding: 0 40px;
  height: 100%;

  ${({ disableTOC }) =>
    !disableTOC &&
    css`
      display: flex;
      justify-content: flex-start;
      align-items: flex-start;
      position: relative;

      @media (max-width: 1200px) {
        flex-direction: column;
      }
    `}

  @media (max-width: 780px) {
    padding: 24px 24px 48px 24px;
  }
`;

export const Children = styled.div`
  width: 100%;
  min-width: 75%;
  max-width: 75%;

  @media (max-width: 1200px) {
    min-width: 100%;
    max-width: 100%;
  }

  ${({ hasTitle }) => !hasTitle && 'padding-top: 40px'};
`;

export const Wrapper = styled.div`
  padding-left: 280px;
  transition: transform 0.5s;

  @media (max-width: 780px) {
    padding-left: 0;
    transform: translate3d(
      ${({ isMenuOpen }) => (isMenuOpen ? '240px' : '0')},
      0,
      0
    );
  }
`;

export const Title = styled.h1`
  padding: 40px 0 0 40px;

  @media (max-width: 780px) {
    padding: 24px 0 0 24px;
  }
`;
