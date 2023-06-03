import { User } from "../lib/goldapps/types";
import { ItemAdded } from "./ItemAdded";

interface Props {
  user: User;
}

export const UserAdditionRow = ({ user }: Props) => {
  return (
    <>
      <td>
        <ItemAdded>{user.cid}</ItemAdded>
      </td>
      <td>
        <ItemAdded>
          {user.first_name} &apos;{user.nick}&apos; {user.second_name}
        </ItemAdded>
      </td>
      <td>
        <ItemAdded>{user.mail}</ItemAdded>
      </td>
      <td>User Addition</td>
    </>
  );
};
