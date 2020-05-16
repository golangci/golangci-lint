/** @jsx jsx */
import { css, jsx } from "@emotion/core";

const ResponsiveContainer = ({ children }) => (
  <div
    css={css`
      max-width: 100%;
      overflow-x: auto;
    `}
  >
    {children}
  </div>
);

export default ResponsiveContainer;
