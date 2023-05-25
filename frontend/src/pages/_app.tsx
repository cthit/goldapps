import type { AppProps } from "next/app";
import "../global.css";

export default function CustomApp({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />;
}
