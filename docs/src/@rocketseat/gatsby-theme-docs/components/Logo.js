import React from "react";
import { css, keyframes } from "@emotion/core";
import GolangciSvg from "./logo.svg";

const grow = keyframes`
0% {
  transform: scale(0.25);
}

35% {
  transform: scale(1.0);
}

70% {
  transform: scale(0.85);
}

100% {
  transform: scale(1);
  opacity: 1;
}
`;

const moveInDown = keyframes`
  0% {
    transform: translate3d(0, -300px, 0);
  }

  60% {
    transform: translate3d(0, 13px, 0);
  }

  80% {
    transform: translate3d(0, -9px, 0);
  }

  100% {
    transform: translate3d(0, 0, 0);
    opacity: 1;
  }
}
`;

const moveInRightShield = keyframes`
0% {
  transform: translate3d(-100px, 46px, 0);
}

50% {
  transform: translate3d(48px, 46px, 0);
}

80% {
  transform: translate3d(0px, 46px, 0);
}

100% {
  transform: translate3d(12px, 46px, 0);
  opacity: 1;
}
`;

const centerAndFade = css`
  opacity: 0;
  transform-origin: 50%;
`;

const svgCss = css`
  width: 100px;
  height: 100px;
  padding: 0.5em;
  &:hover {
    #logo__go__circle {
      animation: ${grow} 0.5s ease-out forwards;
      ${centerAndFade};
    }
    #logo__go__body {
      animation: ${moveInDown} 1s ease-out forwards;
      ${centerAndFade};
    }
    #logo__go__shield {
      animation: ${moveInRightShield} 1s ease-out forwards;
      ${centerAndFade};
    }
  }
`;

const Logo = () => (
  <GolangciSvg x={0} height="80%" viewBox="0 0 100 100" css={svgCss} />
);
export default Logo;
