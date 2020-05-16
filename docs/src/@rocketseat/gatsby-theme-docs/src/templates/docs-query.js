import { graphql } from "gatsby";
import DocsComponent from "../components/Docs-wrapper";

export default DocsComponent;

export const query = graphql`
  query($slug: String!) {
    mdx(fields: { slug: { eq: $slug } }) {
      id
      excerpt(pruneLength: 160)
      fields {
        slug
      }
      frontmatter {
        title
        description
        image
        disableTableOfContents
      }
      body
      headings {
        depth
        value
      }
    }
  }
`;
