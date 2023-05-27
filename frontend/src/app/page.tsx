import { Suggestion } from "../lib/goldapps/types";
import { Suggestions } from "../components/Suggestions";

async function getSuggestions() {
  const mockSuggestions: Suggestion[] = [
    {
      type: "AddUser",
      user: {
        cid: "cid",
        first_name: "first_name",
        second_name: "second_name",
        nick: "nick",
        mail: "mail@chalmers.it",
      },
    },
    {
      type: "DeleteUser",
      user: {
        cid: "cid",
        first_name: "first_name",
        second_name: "second_name",
        nick: "nick",
        mail: "mail@chalmers.it",
      },
    },
    {
      type: "ChangeUser",
      before: {
        cid: "cid",
        first_name: "first_name",
        second_name: "second_name",
        nick: "nick",
        mail: "mail@chalmers.it",
      },
      after: {
        cid: "newcid",
        first_name: "newfirst_name",
        second_name: "newsecond_name",
        nick: "newnick",
        mail: "newmail@chalmers.it",
      },
    },
    {
      type: "AddGroup",
      group: {
        email: "email@chalmers.it",
        expendable: false,
        type: "SOCIETY",
        members: ["member1@chalmers.it"],
      },
    },
    {
      type: "DeleteGroup",
      group: {
        email: "email@chalmers.it",
        expendable: false,
        type: "SOCIETY",
        members: ["member1@chalmers.it"],
      },
    },
    {
      type: "ChangeGroup",
      before: {
        email: "email@chalmers.it",
        expendable: false,
        type: "SOCIETY",
        members: ["member1@chalmers.it", "member2@chalmers.it"],
      },
      after: {
        email: "newemail@chalmers.it",
        expendable: false,
        type: "COMMITTEE",
        members: ["member2@chalmers.it", "member3@chalmers.it"],
      },
    },
  ];

  return mockSuggestions;
}

export default async function IndexPage() {
  const suggestions = await getSuggestions();
  return <Suggestions suggestions={suggestions} />;
}
