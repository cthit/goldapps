import { Group } from "../lib/goldapps/types";
import { arrayDiff, getIdFromEmail } from "../utils";
import { ItemUnchanged } from "./ItemUnchanged";
import { ItemAdded } from "./ItemAdded";
import { ItemDeleted } from "./ItemDeleted";

interface Props {
  before: Group;
  after: Group;
}

export const GroupUpdateRow = ({ before, after }: Props) => {
  const { additions, deletions, unchanged } = arrayDiff(
    before.members || [],
    after.members || [],
  );

  return (
    <>
      <td>
        <ItemUnchanged>{getIdFromEmail(before.email)}</ItemUnchanged>
      </td>
      <td>
        <ItemUnchanged>{before.email}</ItemUnchanged>
      </td>
      <td>
        {unchanged.map(member => (
          <ItemUnchanged key={member}>{member}</ItemUnchanged>
        ))}
        {deletions.map(member => (
          <ItemDeleted key={member}>{member}</ItemDeleted>
        ))}
        {additions.map(member => (
          <ItemAdded key={member}>{member}</ItemAdded>
        ))}
      </td>
      <td>Group Update</td>
    </>
  );
};
