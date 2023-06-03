import { Group } from "../lib/goldapps/types";
import { getIdFromEmail } from "../utils";
import { ItemDeleted } from "./ItemDeleted";

interface Props {
  group: Group;
}

export const GroupDeletionRow = ({ group }: Props) => {
  return (
    <>
      <td>
        <ItemDeleted>{getIdFromEmail(group.email)}</ItemDeleted>
      </td>
      <td>
        <ItemDeleted>{group.email}</ItemDeleted>
      </td>
      <td>
        {group.members &&
          group.members.map(member => (
            <ItemDeleted key={member}>{member}</ItemDeleted>
          ))}
      </td>
      <td>Group Deletion</td>
    </>
  );
};
