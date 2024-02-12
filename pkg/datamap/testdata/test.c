#include <stdio.h>

typedef struct {
  int a;
  int b;
  struct {
    int nested_a;
    int nested_b;
    struct {
      int deeply_nested_a;
      int deeply_nested_b;
    } deeply_nested;
  } nested;
} test_t;

int main() {
  test_t t;
  t.a = 1;
  t.b = 2;
  t.nested.nested_a = 3;
  t.nested.nested_b = 4;
  t.nested.deeply_nested.deeply_nested_a = 5;
  t.nested.deeply_nested.deeply_nested_b = 6;

  printf("t.a: %d\n", t.a);
  printf("t.b: %d\n", t.b);
  printf("t.nested.nested_a: %d\n", t.nested.nested_a);
  printf("t.nested.nested_b: %d\n", t.nested.nested_b);
  printf("t.nested.deeply_nested.deeply_nested_a: %d\n",
         t.nested.deeply_nested.deeply_nested_a);
  printf("t.nested.deeply_nested.deeply_nested_b: %d\n",
         t.nested.deeply_nested.deeply_nested_b);
  printf("done\n");
  return 0;
}
