export function arrayDiff<T>(oldArray: T[], newArray: T[]) {
  const additions = newArray.filter(item => !oldArray.includes(item));
  const deletions = oldArray.filter(item => !newArray.includes(item));
  const unchanged = oldArray.filter(item => newArray.includes(item));

  return {
    additions,
    deletions,
    unchanged,
  };
}

export const getIdFromEmail = (email: string) =>
  email.slice(0, email.search("@"));
