# Backport TODO

## gosec bugfix backport (v2.23.0 → v2.24.7, cherry-pick bugfixes only)

Fork gosec locally and backport these bugfix commits. Do NOT include new
rules (G707, G408, G123, G122, G121, G120, G119, G118, G113) or expanded
G117 (YAML/XML/TOML).

- [ ] G704: false positive on const URLs
- [ ] G705: false positive for non-HTTP `io.Writer`
- [ ] G705: false positive when guard type cannot be resolved (crash)
- [ ] G120: false positive when `MaxBytesReader` is applied in middleware
- [ ] G120: hang-like blowup in wrapper protection checks
- [ ] G602: regression + false positive for array element access
- [ ] G602: false positive for range-over-array indexing
- [ ] G703, G705: taint analysis false positives
- [ ] G115: false positives and false negatives
- [ ] G407: incorrect detection of fixed IV
- [ ] Panic on float constants in overflow analyzer
- [ ] Panic when scanning multi-module repos from root
- [ ] SSA dependency cycle prevention (work amplification on large codebases)
- [ ] SARIF output: invalid null relationships for rules without CWE mappings
- [ ] Sonar report schema compliance fix

Assess cherry-pick feasibility: taint analyzer bugfixes may be entangled
with new taint-based rules (G707, G118, G119, etc.).

