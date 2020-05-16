import styled from '@emotion/styled';
import { darken } from 'polished';

export const Container = styled.aside`
  width: 20%;
  max-width: 280px;
  min-width: 280px;
  background-color: ${({ theme }) => theme.colors.sidebar.background};

  position: fixed;
  overflow-y: auto;
  left: 0;

  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;

  transition: transform 0.5s;

  height: 100vh;

  nav {
    width: 100%;
    align-self: flex-start;
    margin-bottom: 20px;
    flex: 1;
  }

  footer {
    padding: 24px 0 24px 30px;
    width: 100%;

    p {
      color: ${({ theme }) => theme.colors.sidebar.footer};
      font-size: 12px;
      margin: 0;
    }
  }

  @media (max-width: 780px) {
    max-width: 240px;
    min-width: 240px;
    transform: translate3d(
      ${({ isMenuOpen }) => (isMenuOpen ? '0' : '-100%')},
      0,
      0
    );
  }
`;

export const LogoContainer = styled.div`
  width: 100%;
  height: 100%;
  max-height: 100px;
  min-height: 100px;
  padding: 20px 0;

  a {
    width: 100%;
    height: 100%;
    padding-left: 30px;

    display: flex;
    justify-content: flex-start;
    align-items: center;
  }
`;

export const List = styled.ul`
  list-style: none;
  width: 100%;
  padding-left: 0;

  display: flex;
  justify-content: flex-start;
  align-items: center;
  flex-direction: column;
`;

export const Heading = styled.li`
  padding-left: 30px;
  width: 100%;
  text-transform: uppercase;
  font-size: 13px;
  font-weight: bold;
  margin-top: 20px;
  color: ${({ theme }) => theme.colors.primary};
  letter-spacing: 0.142em;
`;

export const Item = styled.li`
  font-size: 15px;
  width: 100%;
  transition: all 200ms ease-in-out;
  padding: 0 20px;

  a,
  span {
    display: block;
    font-size: 15px;
    color: ${({ theme }) => theme.colors.sidebar.link};
    background-color: ${({ theme }) => theme.colors.sidebar.background};
    padding: 4px 10px;
    margin: 4px 0;
    border-radius: 4px;
    font-weight: normal;

    text-decoration: none;
    width: 100%;
    height: 100%;

    display: flex;
    justify-content: flex-start;
    align-items: center;

    cursor: pointer;
    margin: 0 auto;

    transition: background-color 0.2s, color 0.2s, padding-left 0.2s;

    svg {
      width: 20px;
      height: 20px;
      margin-right: 10px;
    }

    &:not(.active-link):hover {
      padding-left: 20px;
      color: ${({ theme }) => darken('0.2', theme.colors.sidebar.link)};
    }

    &.active-link {
      color: ${({ theme }) => darken('0.2', theme.colors.sidebar.link)};
      background-color: ${({ theme }) => theme.colors.sidebar.itemActive};
    }
  }
`;

export const SubItem = styled(List)`
  margin-top: 5px;
`;
