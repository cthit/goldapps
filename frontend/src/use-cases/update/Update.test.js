import * as assert from "assert";
import { filterData, formatData, formatEntry, getAllIds } from "./Update";

const entries = [
  {
    cid: "cid123",
  },
  {
    email: "cid123@chalmers.it",
    members: ["me", "myself", "I"],
  },
  {
    before: {
      cid: "cid123",
    },
    after: {
      cid: "cid123",
    },
  },
  {
    before: {
      email: "cid123@chalmers.it",
      members: ["me", "myself", "I"],
    },
    after: {
      email: "cid123@chalmers.it",
      members: ["me", "myself", "I"],
    },
  },
];

const changeSuggestion = {
  userChanges: {
    userUpdates: [
      {
        before: {
          cid: "damolldd",
          first_name: "Dav",
          second_name: "Mol",
          nick: "Damol",
          mail: "damolldd@student.chalmers.it",
        },
        after: {
          cid: "damolldd",
          first_name: "Dav",
          second_name: "Mol",
          nick: "Damol",
          mail: "damolldd@mol.se",
        },
      },
    ],
    additions: [
      {
        cid: "jhanhard",
        first_name: "Joyce",
        second_name: "Hanhard",
        nick: "Marängsviss",
        mail: "jhanhard@student.chalmers.it",
      },
    ],
    deletions: [
      {
        cid: "davidm",
        first_name: "David",
        second_name: "Möller",
        nick: "Mölle",
        mail: "david@moller.se",
      },
    ],
  },
  groupChanges: {
    groupUpdates: [
      {
        before: {
          email: "kandidatmiddagen2022@chalmers.it",
          type: "FUNCTIONARIES",
          members: [
            "svanni@chalmers.it",
            "hsparritt@chalmers.it",
            "wbulloch@chalmers.it",
            "flongshaw@chalmers.it",
            "davidm@chalmers.it",
          ],
          aliases: null,
          expendable: false,
        },
        after: {
          email: "kandidatmiddagen2022@chalmers.it",
          type: "FUNCTIONARIES",
          members: [
            "svanni@chalmers.it",
            "hsparritt@chalmers.it",
            "wbulloch@chalmers.it",
            "flongshaw@chalmers.it",
          ],
          aliases: null,
          expendable: false,
        },
      },
      {
        before: {
          email: "prit2022@chalmers.it",
          type: "COMMITTEE",
          members: [
            "mjuorio@chalmers.it",
            "tltwell@chalmers.it",
            "hborgba@chalmers.it",
          ],
          aliases: null,
          expendable: false,
        },
        after: {
          email: "prit2022@chalmers.it",
          type: "COMMITTEE",
          members: [
            "jhanhard@chalmers.it",
            "mjuorio@chalmers.it",
            "tltwell@chalmers.it",
            "hborgba@chalmers.it",
          ],
          aliases: null,
          expendable: false,
        },
      },
    ],
    additions: [
      {
        email: "digit2021@chalmers.it",
        type: "ALUMNI",
        members: [
          "jhanhard@student.chalmers.it",
          "tltwell@student.chalmers.it",
          "mjuorio@student.chalmers.it",
          "hborgba@student.chalmers.it",
        ],
        aliases: null,
        expendable: false,
      },
    ],
    deletions: [
      {
        email: "digit2023@chalmers.it",
        type: "ALUMNI",
        members: [
          "mjuorio@student.chalmers.it",
          "tltwell@student.chalmers.it",
          "hborgba@student.chalmers.it",
        ],
        aliases: null,
        expendable: false,
      },
    ],
  },
};

describe("formatEntry: Should return the same object, but with id", () => {
  it("Format user add/delete", () => {
    const entry = formatEntry(entries[0]);
    assert.equal(entry.id, "cid123");
  });
  it("Format group add/delete", () => {
    const entry = formatEntry(entries[1]);
    assert.equal(entry.id, "cid123");
  });
  it("Format user change", () => {
    const entry = formatEntry(entries[2]);
    assert.equal(entry.id, "cid123");
  });
  it("Format group change", () => {
    const entry = formatEntry(entries[3]);
    assert.equal(entry.id, "cid123");
  });
});

describe("formatData: Should append id to each object", () => {
  const [suggestion, ids] = formatData(changeSuggestion);
  it("User changes", () => {
    assert.equal(
      suggestion.userChanges.userUpdates[0].id,
      suggestion.userChanges.userUpdates[0].before.cid,
    );
    assert.equal(ids.includes(suggestion.userChanges.userUpdates[0].id), true);
    assert.equal(
      suggestion.userChanges.additions[0].id,
      suggestion.userChanges.additions[0].cid,
    );
    assert.equal(ids.includes(suggestion.userChanges.additions[0].id), true);
    assert.equal(
      suggestion.userChanges.deletions[0].id,
      suggestion.userChanges.deletions[0].cid,
    );
    assert.equal(ids.includes(suggestion.userChanges.deletions[0].id), true);
  });

  it("Group changes", () => {
    assert.equal(
      suggestion.groupChanges.groupUpdates[0].id,
      "kandidatmiddagen2022",
    );
    assert.equal(ids.includes("kandidatmiddagen2022"), true);
    assert.equal(suggestion.groupChanges.additions[0].id, "digit2021");
    assert.equal(ids.includes("digit2021"), true);
    assert.equal(suggestion.groupChanges.deletions[0].id, "digit2023");
    assert.equal(ids.includes("digit2023"), true);
  });
});

describe("getAllIds: Should give all ids", () => {
  const ids = getAllIds(changeSuggestion);
  it("All ids are included", () => {
    assert.equal(ids.includes("damolldd"), true);
    assert.equal(ids.includes("jhanhard"), true);
    assert.equal(ids.includes("davidm"), true);
    assert.equal(ids.includes("kandidatmiddagen2022"), true);
    assert.equal(ids.includes("prit2022"), true);
    assert.equal(ids.includes("digit2021"), true);
    assert.equal(ids.includes("digit2023"), true);
  });
  it("No extra ids", () => {
    assert.equal(ids.length, 7);
  });
});

describe("filterData: Should filter entries, given a list of ids", () => {
  const [formattedData] = formatData(changeSuggestion);
  it("Filter updates", () => {
    const filteredData = filterData(JSON.parse(JSON.stringify(formattedData)), [
      "damolldd",
      "kandidatmiddagen2022",
      "prit2022",
    ]);

    assert.equal(filteredData.userChanges.userUpdates.length, 1);
    assert.equal(filteredData.userChanges.additions, null);
    assert.equal(filteredData.userChanges.deletions, null);

    assert.equal(filteredData.groupChanges.groupUpdates.length, 2);
    assert.equal(filteredData.groupChanges.additions, null);
    assert.equal(filteredData.groupChanges.deletions, null);
  });

  it("Filter additions", () => {
    const filteredData = filterData(JSON.parse(JSON.stringify(formattedData)), [
      "jhanhard",
      "digit2021",
    ]);

    assert.equal(filteredData.userChanges.userUpdates, null);
    assert.equal(filteredData.userChanges.additions.length, 1);
    assert.equal(filteredData.userChanges.deletions, null);

    assert.equal(filteredData.groupChanges.groupUpdates, null);
    assert.equal(filteredData.groupChanges.additions.length, 1);
    assert.equal(filteredData.groupChanges.deletions, null);
  });

  it("Filter deletions", () => {
    const filteredData = filterData(JSON.parse(JSON.stringify(formattedData)), [
      "davidm",
      "digit2023",
    ]);

    assert.equal(filteredData.userChanges.userUpdates, null);
    assert.equal(filteredData.userChanges.additions, null);
    assert.equal(filteredData.userChanges.deletions.length, 1);

    assert.equal(filteredData.groupChanges.groupUpdates, null);
    assert.equal(filteredData.groupChanges.additions, null);
    assert.equal(filteredData.groupChanges.deletions.length, 1);
  });
});
