export type ContributorInfo = {
  login: string
  name: string
  avatarUrl: string
  websiteUrl: string
}

export type DataJSON = {
  contributors: ContributorInfo[]
  coreTeam: ContributorInfo[]
}
