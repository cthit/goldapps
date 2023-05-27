import { User } from "../lib/goldapps/types";
import { ItemAdded } from "./ItemAdded";
import { ItemDeleted } from "./ItemDeleted";
import { ItemUnchanged } from "./ItemUnchanged";

interface Props {
  before: User;
  after: User;
}

export const UserUpdateRow = ({ before, after }: Props) => {
  return (
    <>
      <td>
        <ItemUnchanged>{before.cid}</ItemUnchanged>
      </td>
      <td>
        {!userNameEqual(before, after) ? (
          <>
            <ItemDeleted>
              {before.first_name} &apos;{before.nick}&apos; {before.second_name}
            </ItemDeleted>
            <ItemAdded>
              {after.first_name} &apos;{after.nick}&apos; {after.second_name}
            </ItemAdded>
          </>
        ) : (
          <>
            <ItemUnchanged>
              {before.first_name} &apos;{before.nick}&apos; {before.second_name}
            </ItemUnchanged>
          </>
        )}
      </td>
      <td>
        {!userEmailEqual(before, after) ? (
          <>
            <ItemDeleted>{before.mail}</ItemDeleted>
            <ItemAdded>{after.mail}</ItemAdded>
          </>
        ) : (
          <ItemUnchanged>{before.mail}</ItemUnchanged>
        )}
      </td>
      <td>User Update</td>
    </>
  );
};

function userNameEqual(before: User, after: User) {
  return (
    before.first_name === after.first_name &&
    before.nick === after.nick &&
    before.second_name === after.second_name
  );
}

function userEmailEqual(before: User, after: User) {
  return before.mail === after.mail;
}
