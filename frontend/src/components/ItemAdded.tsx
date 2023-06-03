export const ItemAdded = ({ children }: React.PropsWithChildren) => {
  return (
    <div>
      <span className="my-1 rounded bg-green-500 px-1 font-mono font-bold">
        + {children}
      </span>
    </div>
  );
};
