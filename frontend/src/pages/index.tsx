import Head from "next/head";
import Home from "../../components/Home";
import { STRINGS } from "../../components/keys";

export default function home() {
  return (
    <>
      <Head>
        <title>{STRINGS.app.title}</title>
        <meta name="keywords" content={STRINGS.app.metaKeywords} />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Home />
    </>
  );
}
