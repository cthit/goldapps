import * as assert from "assert";
import { getId } from "./utils";

describe("getId: Tests if returns id", () => {
  it("Should give id before @ in email", () => {
    assert.equal(getId("hello@googl.com"), "hello");
    assert.equal(
      getId("this.is.an.email2022@this.domain.com"),
      "this.is.an.email2022",
    );
  });
});
