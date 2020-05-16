import styled from '@emotion/styled';

export const Container = styled.section`
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 48px 0;
  width: 100%;

  @media (max-width: 780px) {
    flex-direction: column;
  }
`;

export const Post = styled.div`
  transition: all 200ms;
  ${({ isLeft }) => !isLeft && 'margin-left: auto;'}

  a {
    display: flex;
    text-decoration: none;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;

    svg {
      width: 25px;
      height: 25px;
      color: #737380;

      ${({ isLeft }) => (isLeft ? 'margin-right: 16px' : 'margin-left: 16px')};
    }

    p {
      letter-spacing: 0.142em;
      text-transform: uppercase;
      font-size: 12px;
      margin: 0;
      color: #999;
    }

    h3 {
      color: #737380;
      border: none;
      margin: 0;
      padding: 0;
      font-weight: normal;
      font-size: 16px;
    }
  }

  &:hover {
    opacity: 0.8;

    a svg {
      opacity: 0.8;
    }
  }

  @media (max-width: 780px) {
    width: 100%;
    ${({ isLeft }) => isLeft && 'margin-bottom: 16px'};

    a {
      justify-content: ${({ isLeft }) => (isLeft ? 'flex-start' : 'flex-end')};
    }
  }
`;
