const path = require(`path`);
const { createFilePath } = require(`gatsby-source-filesystem`);
const fs = require(`fs`);

const {
  normalizeBasePath,
  resolveLink,
} = require(`./src/@rocketseat/gatsby-theme-docs-core/util/url`);
const withDefault = require(`./src/@rocketseat/gatsby-theme-docs-core/util/with-default`);

exports.createPages = ({ graphql, actions: { createPage }, reporter }) => {
  reporter.success(`onCreateDocs`);

  const siteConfig = require(`./src/config/site.js`);
  const { basePath, baseDir, docsPath, githubUrl } = withDefault(siteConfig);

  const docsTemplate = require.resolve(
    `./src/@rocketseat/gatsby-theme-docs/src/templates/docs-query.js`
  );
  const homeTemplate = require.resolve(
    `./src/@rocketseat/gatsby-theme-docs/src/templates/homepage-query.js`
  );

  return graphql(
    `
      {
        files: allFile(filter: { extension: { in: ["md", "mdx"] } }) {
          edges {
            node {
              id
              relativePath
              childMdx {
                frontmatter {
                  next
                }
                fields {
                  slug
                }
              }
            }
          }
        }
        sidebar: allSidebarItems {
          edges {
            node {
              label
              link
              items {
                label
                link
              }
              id
            }
          }
        }
      }
    `
  ).then((result) => {
    if (result.errors) {
      reporter.panic(
        `docs: there was an error loading the docs folder!`,
        result.errors
      );
      return;
    }

    createPage({
      path: basePath,
      component: homeTemplate,
    });

    // Generate prev/next items based on sidebar.yml file
    const sidebar = result.data.sidebar.edges;
    const listOfItems = [];

    sidebar.forEach(({ node: { label, link, items } }) => {
      if (Array.isArray(items)) {
        items.forEach((item) => {
          listOfItems.push({
            label: item.label,
            link: resolveLink(item.link, basePath),
          });
        });
      } else {
        listOfItems.push({
          label,
          link: resolveLink(link, basePath),
        });
      }
    });

    // Generate docs pages
    const docs = result.data.files.edges;
    docs.forEach((doc) => {
      const {
        childMdx: {
          frontmatter: { next },
          fields: { slug },
        },
        relativePath,
      } = doc.node;

      const githubEditUrl =
        githubUrl &&
        `${githubUrl}/tree/master/${baseDir}/${docsPath}/${relativePath}`;

      const currentPageIndex = listOfItems.findIndex(
        (page) => page.link === slug
      );

      const prevItem = listOfItems[currentPageIndex - 1];
      const nextItem = next
        ? listOfItems.find((item) => item.link === next)
        : listOfItems[currentPageIndex + 1];

      createPage({
        path: slug,
        component: docsTemplate,
        context: {
          slug,
          prev: prevItem,
          next: nextItem,
          githubEditUrl,
        },
      });
    });

    reporter.success(`docs pages created`);
  });
};

exports.createSchemaCustomization = ({ actions }) => {
  actions.createTypes(`
    type MdxFrontmatter @dontInfer {
      title: String!
      description: String
      image: String
      disableTableOfContents: Boolean
      next: String
    }
  `);

  actions.createTypes(`
    type SidebarItems implements Node {
      label: String!
      link: String
      items: [SidebarItemsItems]
    }

    type SidebarItemsItems {
      label: String
      link: String
    }
  `);
};

exports.onPreBootstrap = ({ store, reporter }, themeOptions) => {
  const { configPath, docsPath } = withDefault(themeOptions);
  const { program } = store.getState();

  const dirs = [
    path.join(program.directory, configPath),
    path.join(program.directory, docsPath),
  ];

  dirs.forEach((dir) => {
    if (!fs.existsSync(dir)) {
      reporter.success(`docs: intialized the ${dir} directory`);
      fs.mkdirSync(dir);
    }
  });
};

exports.onCreateNode = (
  { node, actions: { createNodeField }, getNode },
  themeOptions
) => {
  if (node.internal.type !== `Mdx`) {
    return;
  }

  const { basePath } = withDefault(themeOptions);

  let value = createFilePath({ node, getNode });
  if (value === "index") value = "";

  createNodeField({
    name: `slug`,
    node,
    value: normalizeBasePath(basePath, value),
  });

  createNodeField({
    name: `id`,
    node,
    value: node.id,
  });
};

/**
[
  {
    "node": {
      "label": "Home",
      "link": "/",
      "items": null,
      "id": "a2913be3-af3c-5fc9-967e-a058e86b20a9"
    }
  },
  {
    "node": {
      "label": "With dropdown",
      "link": null,
      "items": [
        { "label": "My Example", "link": "/my-example" },
        { "label": "Teste 2", "link": "/teste-2" }
      ],
      "id": "c7d9606c-4bda-5097-a0df-53108e9f4efd"
    }
  }
]
*/

// Ler todo o array e salvar em uma objeto chave/valor
/**
 * {
 *    '/': {
 *       prev: null,
 *       next: {
 *          label: 'My example',
 *          link: '/my-example'
 *       }
 *    },
 *    '/my-example': {
 *       prev: {
 *          label: 'Home',
 *          link: '/'
 *       },
 *       next: {
 *          label: 'Teste 2',
 *          link: '/teste-2'
 *       }
 *    },
 *    '/teste-2': {
 *       prev: {
 *          label: 'My example',
 *          link: '/my-example'
 *       },
 *       next: null
 *    }
 * }
 */
