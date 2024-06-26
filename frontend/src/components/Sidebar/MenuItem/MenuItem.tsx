import { NavLink } from "react-router-dom";

import { linkStyles, wrapperStyles } from "./MenuItem.css";

export type MenuItemProps = {
  title: string;
  id: number;
};

export const MenuItem = ({ title, id }: MenuItemProps) => {
  return (
    <div className={wrapperStyles}>
      <NavLink to={`/chat/${id}`} className={linkStyles}>
        {title}
      </NavLink>
    </div>
  );
};
