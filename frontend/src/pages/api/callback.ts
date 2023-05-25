import { NextApiRequest, NextApiResponse } from "next";

const handler = async (req: NextApiRequest, res: NextApiResponse) => {
  try {
    res.redirect(
      "http://localhost:8080/api/authenticate?code=" + req.query.code,
    );
  } catch (e) {
    res.redirect("/unauthorized");
  }
};

export default handler;
