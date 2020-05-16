module.exports = themeOptions => {
  const basePath = themeOptions.basePath || `/`;
  const configPath = themeOptions.configPath || `config`;
  const docsPath = themeOptions.docsPath || `docs`;
  const baseDir = themeOptions.baseDir || ``;
  const { githubUrl } = themeOptions;

  return {
    basePath,
    configPath,
    docsPath,
    baseDir,
    githubUrl,
  };
};
