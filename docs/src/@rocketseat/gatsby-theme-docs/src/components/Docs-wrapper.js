/* eslint-disable react/prop-types */
import React from "react";
import Docs from "./Docs";

export default function Docspage({ data: { mdx }, pageContext }) {
  return <Docs mdx={mdx} pageContext={pageContext} />;
}
