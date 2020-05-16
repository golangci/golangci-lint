function normalizeBasePath(basePath, link) {
  return `/${basePath}/${link}`.replace(/\/\/+/g, `/`);
}

function isExternalUrl(url) {
  return new RegExp('^((https?:)?//)', 'i').test(url);
}

function resolveLink(link, basePath) {
  return isExternalUrl(link) ? link : normalizeBasePath(basePath, link);
}

module.exports = { resolveLink, normalizeBasePath, isExternalUrl };
