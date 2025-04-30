import { Stack } from "expo-router";

export default function AuthLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: "Login" }} />
      <Stack.Screen name="signup" options={{ title: "Sign Up" }} />
      <Stack.Screen name="confirm" options={{ title: "Confirm Email" }} />
    </Stack>
  );
}
