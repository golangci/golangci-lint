function isExternalUrl(url) {
  return new RegExp('^((https?:)?//)', 'i').test(url);
}

module.exports = { isExternalUrl };
