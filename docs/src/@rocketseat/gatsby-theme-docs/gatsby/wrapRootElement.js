/* eslint-disable */
import React from 'react';
import { ThemeProvider } from 'emotion-theming';

import defaultTheme from '../src/styles/theme';
import GlobalStyle from '../src/styles/global';

export function wrapRootElement({ element }) {
  return (
    <ThemeProvider theme={defaultTheme}>
      <>
        <GlobalStyle />
        {element}
      </>
    </ThemeProvider>
  );
}
