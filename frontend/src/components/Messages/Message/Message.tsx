import { formatDistance } from "date-fns";
import { useState } from "react";

import { Button } from "@/components/Button/Button";
import { Icon } from "@/components/Icon/Icon";
import { TaskStatus, TaskType } from "@/generated/graphql";

import {
  avatarStyles,
  contentStyles,
  iconStyles,
  messageStyles,
  outputStyles,
  rightColumnStyles,
  timeStyles,
  wrapperStyles,
} from "./Message.css";

type MessageProps = {
  message: string;
  time: Date;
  type: TaskType;
  status: TaskStatus;
  output: string;
};

export const Message = ({
  time,
  message,
  type,
  status,
  output,
}: MessageProps) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleExpand = () => {
    setIsExpanded((prev) => !prev);
  };

  return (
    <div className={wrapperStyles}>
      <img
        src="https://via.placeholder.com/40"
        alt="avatar"
        className={avatarStyles}
        width="40"
        height="40"
      />
      <div className={rightColumnStyles}>
        <div className={timeStyles}>
          {formatDistance(new Date(time), new Date(), { addSuffix: true })}
        </div>
        <div
          className={
            type === TaskType.Input
              ? messageStyles.Input
              : status !== TaskStatus.Failed
                ? messageStyles.Regular
                : messageStyles.Failed
          }
          onClick={toggleExpand}
        >
          <div className={contentStyles}>
            <span
              className={
                status !== TaskStatus.Failed
                  ? iconStyles.Regular
                  : iconStyles.Failed
              }
            >
              {getIcon(type)}
            </span>
            <div>{message}</div>
          </div>
          {status === TaskStatus.InProgress && (
            <Button size="small" hierarchy="danger">
              Stop
            </Button>
          )}
        </div>
        {isExpanded && <div className={outputStyles}>{output}</div>}
      </div>
    </div>
  );
};

const getIcon = (type: TaskType) => {
  let icon = null;

  switch (type) {
    case TaskType.Browser:
      icon = <Icon.Browser />;
      break;
    case TaskType.Terminal:
      icon = <Icon.Terminal />;
      break;
    case TaskType.Code:
      icon = <Icon.Code />;
      break;
    case TaskType.Ask:
      icon = <Icon.MessageQuestion />;
      break;
    case TaskType.Done:
      icon = <Icon.CheckCircle />;
      break;
  }

  return icon;
};
