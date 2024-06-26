// This terminal is a combination of the following packages:
// https://gist.github.com/mastersign/90d0ab06f040092e4ca27a3b59820cb9
// https://github.com/reubenmorgan/xterm-react/blob/6c8bb143387a6abc35ff54a3e099c46e5be8819c/src/Xterm.tsx
import React, { useEffect, useRef } from "react";
import { ITerminalAddon, ITerminalOptions, Terminal as XTerminal } from "xterm";
import { AttachAddon } from "xterm-addon-attach";
import { CanvasAddon } from "xterm-addon-canvas";
import { FitAddon } from "xterm-addon-fit";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { WebLinksAddon } from "xterm-addon-web-links";
import { WebglAddon } from "xterm-addon-webgl";
import { Broadcast } from "xterm-theme";
import "xterm/css/xterm.css";

const isWebGl2Supported = !!document
  .createElement("canvas")
  .getContext("webgl2");

function useBind(
  termRef: React.RefObject<XTerminal>,
  handler: any,
  eventName:
    | "onBell"
    | "onBinary"
    | "onCursorMove"
    | "onData"
    | "onKey"
    | "onLineFeed"
    | "onRender"
    | "onResize"
    | "onScroll"
    | "onSelectionChange"
    | "onTitleChange"
    | "onWriteParsed",
) {
  useEffect(() => {
    if (!termRef.current || typeof handler !== "function") return;
    const term = termRef.current;
    const eventBinding = term[eventName](handler);
    return () => {
      if (!eventBinding) return;
      eventBinding.dispose();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [handler]);
}

type XTermProps = {
  customKeyEventHandler?(event: KeyboardEvent): boolean;
  className?: string;
  id?: string;
  onBell?: () => void;
  onBinary?: (data: string) => void;
  onCursorMove?: () => void;
  onData?: (data: string) => void;
  onDispose?: (term: XTerminal) => void;
  onInit?: (term: XTerminal) => void;
  onKey?: (key: { key: string; domEvent: KeyboardEvent }) => void;
  onLineFeed?: () => void;
  onRender?: () => void;
  onResize?: (cols: number, rows: number) => void;
  onScroll?: (ydisp: number) => void;
  onSelectionChange?: () => void;
  onTitleChange?: (title: string) => void;
  onWriteParsed?: (data: string) => void;
  options?: ITerminalOptions;
};

const addons: ITerminalAddon[] = [
  new Unicode11Addon(),
  new CanvasAddon(),
  isWebGl2Supported ? new WebglAddon() : new WebLinksAddon(),
];

export const Terminal = ({
  id,
  className,
  onBell,
  onBinary,
  onCursorMove,
  onData,
  onKey,
  onLineFeed,
  onRender,
  onResize,
  onScroll,
  onSelectionChange,
  onTitleChange,
  onWriteParsed,
  customKeyEventHandler,
  onInit,
}: XTermProps) => {
  const divRef = useRef<HTMLDivElement | null>(null);
  const xtermRef = useRef<XTerminal | null>(null);
  const connectedRefSocket = useRef<WebSocket>();

  useEffect(() => {
    if (!xtermRef.current) return;

    xtermRef.current.clear();

    if (connectedRefSocket.current) {
      console.log(`Closing connection to the terminal`);
      connectedRefSocket.current.close();
    }

    if (id) {
      const socket = new WebSocket(
        "ws://" + import.meta.env.VITE_API_URL + "/terminal/" + id,
      );
      xtermRef.current.loadAddon(new AttachAddon(socket));
      connectedRefSocket.current = socket;
      console.log(`Connected to terminal #${id}`);
    }
  }, [id]);

  useEffect(() => {
    if (!divRef.current || xtermRef.current) return;
    const xterm = new XTerminal({
      convertEol: true,
      allowProposedApi: true,
      theme: Broadcast,
    });

    // Load addons if the prop exists.
    addons.forEach((addon) => {
      xterm.loadAddon(addon);
    });

    const fitAddon = new FitAddon();
    xterm.loadAddon(fitAddon);

    // Add Custom Key Event Handler if provided
    if (customKeyEventHandler) {
      xterm.attachCustomKeyEventHandler(customKeyEventHandler);
    }

    xtermRef.current = xterm;
    xterm.open(divRef.current);
    fitAddon.fit();
  }, [id]);

  useBind(xtermRef, onBell, "onBell");
  useBind(xtermRef, onBinary, "onBinary");
  useBind(xtermRef, onCursorMove, "onCursorMove");
  useBind(xtermRef, onData, "onData");
  useBind(xtermRef, onKey, "onKey");
  useBind(xtermRef, onLineFeed, "onLineFeed");
  useBind(xtermRef, onRender, "onRender");
  useBind(xtermRef, onResize, "onResize");
  useBind(xtermRef, onScroll, "onScroll");
  useBind(xtermRef, onSelectionChange, "onSelectionChange");
  useBind(xtermRef, onTitleChange, "onTitleChange");
  useBind(xtermRef, onWriteParsed, "onWriteParsed");

  useEffect(
    () => {
      if (!xtermRef.current) return;
      if (typeof onInit !== "function") return;
      onInit(xtermRef.current);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [xtermRef.current],
  );

  return <div id={id} className={className} ref={divRef} />;
};
