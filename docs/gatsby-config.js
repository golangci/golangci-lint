const withDefault = require(`./src/@rocketseat/gatsby-theme-docs-core/util/with-default`);

const domain = `golangci-lint.run`;
const siteUrl = `https://${domain}`;

const siteConfig = require(`./src/config/site.js`);
const { basePath, configPath, docsPath } = withDefault(siteConfig);

module.exports = {
  siteMetadata: {
    siteTitle: `golangci-lint`,
    defaultTitle: ``,
    siteTitleShort: `golangci-lint`,
    siteDescription: `Fast Go linters runner golangci-lint.`,
    siteUrl,
    siteAuthor: `@golangci`,
    siteImage: `/logo.png`,
    siteLanguage: `en`,
    themeColor: `#7159c1`,
    basePath,
    footer: `© ${new Date().getFullYear()}`,
  },
  plugins: [
    `gatsby-alias-imports`,
    `gatsby-transformer-sharp`,
    `gatsby-plugin-sharp`,
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        name: `docs`,
        path: docsPath,
      },
    },
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        name: `config`,
        path: configPath,
      },
    },
    {
      resolve: `gatsby-transformer-yaml`,
      options: {
        typeName: `SidebarItems`,
      },
    },
    {
      resolve: `gatsby-plugin-mdx`,
      options: {
        extensions: [`.mdx`, `.md`],
        gatsbyRemarkPlugins: [
          `gatsby-remark-autolink-headers`,
          `gatsby-remark-external-links`,
          {
            resolve: `gatsby-remark-images`,
            options: {
              maxWidth: 960,
              withWebp: true,
              linkImagesToOriginal: false,
            },
          },
          `gatsby-remark-responsive-iframe`,
          `gatsby-remark-copy-linked-files`,
          `gatsby-remark-mermaid`,
        ],
        plugins: [
          `gatsby-remark-autolink-headers`,
          `gatsby-remark-external-links`,
          `gatsby-remark-images`,
          `gatsby-remark-mermaid`,
        ],
      },
    },
    {
      resolve: `gatsby-plugin-manifest`,
      options: {
        name: `golangci-lint website`,
        short_name: `golangci-lint`,
        start_url: `/`,
        background_color: `#ffffff`,
        display: `standalone`,
        icon: `static/logo.png`,
      },
    },
    `gatsby-plugin-sitemap`,
    {
      resolve: `gatsby-plugin-google-analytics`,
      options: {
        trackingId: `UA-48413061-13`,
        siteSpeedSampleRate: 100,
      },
    },
    {
      resolve: `gatsby-plugin-canonical-urls`,
      options: {
        siteUrl: siteUrl,
      },
    },
    `gatsby-plugin-offline`,
    {
      resolve: "gatsby-plugin-react-svg",
      options: {
        rule: {
          include: /logo\.svg$/,
        },
      },
    },
    `gatsby-plugin-netlify`,
    `gatsby-plugin-netlify-cache`,
    `gatsby-plugin-catch-links`,
    `gatsby-plugin-emotion`,
    `gatsby-plugin-react-helmet`,
    {
      resolve: `gatsby-plugin-robots-txt`,
      options: {
        env: {
          development: {
            host: siteUrl,
            policy: [{ userAgent: "*", disallow: ["/"] }],
          },
          production: {
            host: siteUrl,
            policy: [{ userAgent: "*", disallow: ["/page-data/"] }],
          },
        },
      },
    },
  ],
};
