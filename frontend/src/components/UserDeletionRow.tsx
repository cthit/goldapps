import { User } from "../lib/goldapps/types";
import { ItemDeleted } from "./ItemDeleted";

interface Props {
  user: User;
}

export const UserDeletionRow = ({ user }: Props) => {
  return (
    <>
      <td>
        <ItemDeleted>{user.cid}</ItemDeleted>
      </td>
      <td>
        <ItemDeleted>
          {user.first_name} &apos;{user.nick}&apos; {user.second_name}
        </ItemDeleted>
      </td>
      <td>
        <ItemDeleted>{user.mail}</ItemDeleted>
      </td>
      <td>User Deletion</td>
    </>
  );
};
