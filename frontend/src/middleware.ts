import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const GOLDAPPS_URL = process.env.GOLDAPPS_URL || "http://localhost:8080";
  return NextResponse.rewrite(GOLDAPPS_URL + request.nextUrl.pathname);
}

export const config = {
  matcher: "/api/:path*",
};
