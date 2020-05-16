const siteUrl = `https://golangci-lint.run`;

module.exports = {
  siteMetadata: {
    siteTitle: `golangci-lint`,
    defaultTitle: ``,
    siteTitleShort: `golangci-lint`,
    siteDescription: `Fast Go linters runner golangci-lint.`,
    siteUrl: siteUrl,
    siteAuthor: `@golangci`,
    siteImage: `/logo.png`,
    siteLanguage: `en`,
    themeColor: `#7159c1`,
    basePath: `/`,
    footer: `© ${new Date().getFullYear()}`,
  },
  plugins: [
    `gatsby-alias-imports`,
    {
      resolve: `@rocketseat/gatsby-theme-docs`,
      options: {
        configPath: `src/config`,
        docsPath: `src/docs`,
        githubUrl: `https://github.com/golangci/golangci-lint`,
        baseDir: `docs`,
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
        trackingId: null, // TODO
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
    {
      resolve: `gatsby-transformer-remark`,
      options: {
        plugins: [`gatsby-remark-external-links`],
      },
    },
    `gatsby-plugin-netlify`,
    `gatsby-plugin-netlify-cache`,
  ],
};
