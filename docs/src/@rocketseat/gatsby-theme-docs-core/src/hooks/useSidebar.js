import { graphql, useStaticQuery } from 'gatsby';
import { resolveLink } from '../../util/url';

export function useSidebar() {
  const data = useStaticQuery(graphql`
    {
      allSidebarItems {
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
      site {
        siteMetadata {
          basePath
        }
      }
    }
  `);

  const { basePath } = data.site.siteMetadata;

  const {
    allSidebarItems: { edges },
  } = data;

  if (basePath) {
    const normalizedSidebar = edges.map(
      ({ node: { label, link, items, id } }) => {
        if (Array.isArray(items)) {
          items = items.map(item => ({
            label: item.label,
            link: resolveLink(item.link, basePath),
          }));
        }

        return {
          node: {
            id,
            label,
            items,
            link: resolveLink(link, basePath),
          },
        };
      },
    );

    return normalizedSidebar;
  }

  return edges;
}
