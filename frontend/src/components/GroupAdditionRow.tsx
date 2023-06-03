import { Group } from "../lib/goldapps/types";
import { getIdFromEmail } from "../utils";
import { ItemAdded } from "./ItemAdded";

interface Props {
  group: Group;
}

export const GroupAdditionRow = ({ group }: Props) => {
  return (
    <>
      <td>
        <ItemAdded>{getIdFromEmail(group.email)}</ItemAdded>
      </td>
      <td>
        <ItemAdded>{group.email}</ItemAdded>
      </td>
      <td>
        {group.members &&
          group.members.map(members => (
            <ItemAdded key={members}>{members}</ItemAdded>
          ))}
      </td>
      <td>Group Addition</td>
    </>
  );
};
