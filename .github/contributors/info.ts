export type ContributorInfo = {
  login: string
  name: string
  avatarUrl: string
  websiteUrl: string
  isTeamMember: boolean
}

export type DataJSON = {
  contributors: ContributorInfo[]
  coreTeam: ContributorInfo[]
}
