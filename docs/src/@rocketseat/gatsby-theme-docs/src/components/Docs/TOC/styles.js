import styled from '@emotion/styled';
import { darken } from 'polished';

export const Container = styled.div`
  position: sticky;
  order: 2;
  margin-left: 3rem;
  max-width: 18rem;
  top: 4rem;
  max-height: calc(100vh - 4rem - 2.5rem - 3rem - 3rem);
  overflow: auto;
  width: 100%;
  max-width: 25%;

  h2 {
    color: #737380;
    text-transform: uppercase;
    font-size: 13px;
    font-weight: bold;
    letter-spacing: 0.142em;
    margin-top: 0rem;
    border: none;
    margin: 0 0 24px 0;
  }

  nav ul {
    padding-left: 0;
    margin-bottom: 0;
    list-style: none;

    li {
      margin-bottom: 12px;
      line-height: 1.1;

      a {
        font-size: 13px;
        color: #999999;
        font-weight: 400;
        text-decoration: none;
        transition: all 0.2s;

        &:hover {
          color: ${({ theme }) => darken('0.2', theme.colors.sidebar.link)};
        }
      }
    }
  }

  @media (max-width: 1200px) {
    position: relative;
    top: auto;
    order: 0;
    max-width: 100%;
    margin: 0 0 24px 0;
    padding-bottom: 16px;
    border-bottom: 1px solid rgba(120, 117, 122, 0.2);
  }
`;
