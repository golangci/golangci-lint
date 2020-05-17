import { danger } from "danger";

const comment = (username: string) => `
Hey, @${username} — we just merged your PR to \`golangci-lint\`! 🔥🚀

\`golangci-lint\` is built by awesome people like you. Let us say “thanks”: **we just invited you to join the [GolangCI](https://github.com/golangci) organization on GitHub.**
This will add you to our team of maintainers. Accept the invite by visiting [this link](https://github.com/orgs/golangci/invitation).

By joining the team, you’ll be able to label issues, review pull requests, and merge approved pull requests.
More information about contributing is [here](https://golangci-lint.run/contributing/quick-start/).

Thanks again!
`;

const teamId = `3840765`;

export const inviteCollaborator = async () => {
  const gh = danger.github;
  const api = gh.api;

  // Details about the repo.
  const owner = gh.thisPR.owner;
  const repo = gh.thisPR.repo;
  const number = gh.thisPR.number;

  // Details about the collaborator.
  const username = gh.pr.user.login;

  // Check whether or not we’ve already invited this contributor.
  try {
    const inviteCheck = (await api.teams.getMembership({
      team_id: teamId,
      username,
    } as any)) as any;
    const isInvited = inviteCheck.headers.status !== "404";

    // If we’ve already invited them, don’t spam them with more messages.
    if (isInvited) {
      console.log(
        `@${username} has already been invited to this org. Doing nothing.`
      );
      return;
    }
  } catch (err) {
    console.info(
      `Error checking membership of ${username} in team ${teamId}: ${err.stack}`
    );
    // If the user hasn’t been invited, the invite check throws an error.
  }

  try {
    const invite = await api.teams.addOrUpdateMembership({
      team_id: teamId,
      username,
    } as any);

    if (invite.data.state === "active") {
      console.log(
        `@${username} is already a ${invite.data.role} for this team.`
      );
    } else {
      console.log(`We’ve invited @${username} to join this team.`);
    }
  } catch (err) {
    console.log("Something went wrong.");
    console.log(err);
    return;
  }

  // For new contributors, roll out the welcome wagon!
  await api.issues.createComment({
    owner,
    repo,
    number,
    body: comment(username),
  });
};

export default async () => {
  await inviteCollaborator();
};
