export const ItemDeleted = ({ children }: React.PropsWithChildren) => {
  return (
    <div>
      <span className="my-1 rounded bg-red-500 font-mono font-bold">
        - {children}
      </span>
    </div>
  );
};
