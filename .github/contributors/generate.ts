import nyc from "name-your-contributors"
import graphql from "@octokit/graphql"
import * as fs from "fs"
import { ContributorInfo, DataJSON } from "./info"

type WeightedContributor = {
  login: string
  weight: number
}

type Contribution = {
  login: string
  count: number
}

const buildWeights = (contributorStats: any[]): Map<string, number> => {
  const loginToTotalWeight: Map<string, number> = new Map<string, number>()

  const addContributions = (weight: number, contributions: Contribution[], onlyExisting: boolean) => {
    for (const contr of contributions) {
      if (onlyExisting && !loginToTotalWeight.has(contr.login)) {
        continue
      }
      loginToTotalWeight.set(contr.login, (loginToTotalWeight.get(contr.login) || 0) + contr.count * weight)
    }
  }

  for (const stat of contributorStats) {
    // totally every pull or commit should account as 10
    addContributions(5, stat.prCreators, false)
    addContributions(5, stat.commitAuthors, false) // some commits are out pull requests
  }

  for (const stat of contributorStats) {
    addContributions(2, stat.prCommentators, true)
    addContributions(2, stat.issueCreators, true)
    addContributions(2, stat.issueCommentators, true)
    addContributions(2, stat.reviewers, true)
    addContributions(0.3, stat.reactors, true)
  }

  return loginToTotalWeight
}

const buildContributorInfo = async (contributors: WeightedContributor[]): Promise<ContributorInfo[]> => {
  const query = `{
    ${contributors.map((c, i) => `user${i}: user(login: "${c.login}") {...UserFragment}`).join(`\n`)}
  }
    fragment UserFragment on User {
        login
        name
        websiteUrl
        avatarUrl
    }`

  try {
    const resp = await graphql.graphql(query, {
      headers: {
        authorization: `token ${process.env.GITHUB_TOKEN}`,
      },
    })

    return contributors.map((_, i) => resp[`user${i}`]).filter((v) => v)
  } catch (err) {
    if (err.errors && err.data) {
      console.warn(`github errors:`, err.errors)
      return contributors.map((_, i) => err.data[`user${i}`]).filter((v) => v)
    }
    throw err
  }
}

const buildCoreTeamInfo = async (): Promise<ContributorInfo[]> => {
  const query = `{
        organization(login:"golangci"){
          team(slug:"core-team"){
            members {
              nodes {
                login
                name
                websiteUrl
                avatarUrl
              }
            }
          }
        }
      }`
  const resp = await graphql.graphql(query, {
    headers: {
      authorization: `token ${process.env.GITHUB_TOKEN}`,
    },
  })

  return resp.organization.team.members.nodes
}

const main = async () => {
  try {
    const repos = [`golangci-lint`, `golangci-lint-action`]
    const promises = repos.map((repo) =>
      nyc.repoContributors({
        token: process.env.GITHUB_TOKEN,
        user: `golangci`,
        repo,
        before: new Date(),
        after: new Date(0),
        commits: true,
        reactions: true,
      })
    )
    const contributorStats = await Promise.all(promises)
    console.info(`Got stat for repos ${repos}`)
    const loginToWeight = buildWeights(contributorStats)
    const weightedContributors: WeightedContributor[] = []
    loginToWeight.forEach((weight, login) => weightedContributors.push({ login, weight }))

    weightedContributors.sort((a, b) => b.weight - a.weight)
    const coreTeamInfo = await buildCoreTeamInfo()
    const contributorsInfo = await buildContributorInfo(weightedContributors)
    const exclude: any = {
      golangcidev: true,
      CLAassistant: true,
      renovate: true,
      fossabot: true,
      golangcibot: true,
      kortschak: true,
      "golanci-releaser": true,
    }

    const res: DataJSON = {
      contributors: contributorsInfo.filter((c) => !exclude[c.login] && !coreTeamInfo.find((ct) => ct.login === c.login)),
      coreTeam: coreTeamInfo.sort((a, b) => (loginToWeight.get(b.login) || 0) - (loginToWeight.get(a.login) || 0)),
    }
    console.info(res)
    fs.writeFileSync("contributors.json", JSON.stringify(res, null, 2))
  } catch (err) {
    console.error(`Failed to get repo contributors`, err)
    process.exit(1)
  }
  console.info(`Success`)
}

main()
