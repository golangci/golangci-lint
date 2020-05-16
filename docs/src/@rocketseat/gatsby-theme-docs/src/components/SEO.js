import React from 'react';
import PropTypes from 'prop-types';
import Helmet from 'react-helmet';
import urljoin from 'url-join';
import { useStaticQuery, graphql } from 'gatsby';

export default function SEO({ description, title, slug, image, children }) {
  const { site } = useStaticQuery(
    graphql`
      query {
        site {
          siteMetadata {
            defaultTitle
            siteTitleShort
            siteTitle
            siteImage
            siteDescription
            siteLanguage
            siteUrl
            siteAuthor
          }
        }
      }
    `,
  );

  const {
    siteTitle,
    siteTitleShort,
    siteUrl,
    defaultTitle,
    siteImage,
    siteDescription,
    siteLanguage,
    siteAuthor,
    siteIcon,
  } = site.siteMetadata;

  const metaTitle = title ? `${title} | ${siteTitle}` : defaultTitle;
  const metaUrl = urljoin(siteUrl, slug);
  const metaImage = urljoin(siteUrl, image || siteImage);
  const metaDescription = description || siteDescription;

  const schemaOrgJSONLD = [
    {
      '@context': 'http://schema.org',
      '@type': 'WebSite',
      url: metaUrl,
      name: title,
      alternateName: siteTitleShort,
    },
  ];

  return (
    <Helmet
      htmlAttributes={{
        lang: siteLanguage,
      }}
      title={metaTitle}
    >
      {siteIcon && <link rel="icon" href={siteIcon} />}
      <meta name="description" content={metaDescription} />
      <meta name="image" content={metaImage} />

      <meta httpEquiv="x-ua-compatible" content="IE=edge,chrome=1" />
      <meta name="MobileOptimized" content="320" />
      <meta name="HandheldFriendly" content="True" />
      <meta name="google" content="notranslate" />
      <meta name="referrer" content="no-referrer-when-downgrade" />

      <meta property="og:url" content={metaUrl} />
      <meta property="og:type" content="website" />
      <meta property="og:title" content={metaTitle} />
      <meta property="og:description" content={metaDescription} />
      <meta property="og:locale" content={siteLanguage} />
      <meta property="og:site_name" content={siteTitle} />
      <meta property="og:image" content={metaImage} />
      <meta property="og:image:secure_url" content={metaImage} />
      <meta property="og:image:alt" content="Banner" />
      <meta property="og:image:type" content="image/png" />
      <meta property="og:image:width" content="1200" />
      <meta property="og:image:height" content="630" />

      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={metaTitle} />
      <meta name="twitter:site" content={siteAuthor} />
      <meta name="twitter:creator" content={siteAuthor} />
      <meta name="twitter:image" content={metaImage} />
      <meta name="twitter:image:src" content={metaImage} />
      <meta name="twitter:image:alt" content="Banner" />
      <meta name="twitter:image:width" content="1200" />
      <meta name="twitter:image:height" content="630" />

      <script type="application/ld+json">
        {JSON.stringify(schemaOrgJSONLD)}
      </script>
      {children}
    </Helmet>
  );
}

SEO.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
  slug: PropTypes.string,
  image: PropTypes.string,
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.element),
    PropTypes.node,
  ]),
};

SEO.defaultProps = {
  title: '',
  description: '',
  slug: '',
  image: '',
  children: '',
};
