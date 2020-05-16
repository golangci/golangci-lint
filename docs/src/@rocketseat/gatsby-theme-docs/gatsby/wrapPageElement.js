/* eslint-disable */
import React from 'react';
import { MDXProvider } from '@mdx-js/react';

import Code from '../src/components/Code';

const components = {
  code: Code,
  inlineCode: props => <code className="inline-code" {...props} />,
};

export function wrapPageElement({ element }) {
  return <MDXProvider components={components}>{element}</MDXProvider>;
}
