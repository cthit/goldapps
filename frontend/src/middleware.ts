import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { createGoldappsServerClient } from "./lib/goldapps/client-server";
import { checkLogin } from "./lib/goldapps/auth";

export async function middleware(request: NextRequest) {

  const GOLDAPPS_URL = process.env.GOLDAPPS_URL || "http://localhost:8080";
  if (request.nextUrl.pathname === "/") {
    const loginStatus = await checkLogin();
    if (loginStatus?.data){
      const redirectResponse = NextResponse.redirect(loginStatus.data);

      if (loginStatus.cookie)
        redirectResponse.cookies.set(loginStatus.cookie);

      return redirectResponse;
    }
  }

  if (request.nextUrl.pathname.startsWith("/api/")) {
    return NextResponse.rewrite(
      GOLDAPPS_URL + request.nextUrl.pathname + request.nextUrl.search,
    );
  }
  return NextResponse.next();
}

//export const config = {
//  matcher: "/api/:path*",
//};
