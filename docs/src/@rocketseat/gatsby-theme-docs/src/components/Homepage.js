import React from 'react';
import Index from '../text/index.mdx';

import Layout from './Layout';
import SEO from './SEO';

export default function Home() {
  return (
    <Layout>
      <SEO />
      <Index />
    </Layout>
  );
}
