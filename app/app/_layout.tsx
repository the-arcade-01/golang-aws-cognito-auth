import { Stack } from "expo-router";
import { useEffect, useState } from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function RootLayout() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const token = await AsyncStorage.getItem("token");
        console.log("Auth token: ", token);
        setIsAuthenticated(!!token);
      } catch (error) {
        console.error("Error checking auth:", error);
        setIsAuthenticated(false);
      }
    };
    checkAuth();
  }, []);

  if (isAuthenticated === null) {
    return null;
  }

  return (
    <Stack screenOptions={{ headerShown: false }}>
      {/* {isAuthenticated ? (
        <Stack.Screen name="(tabs)" />
      ) : (
        <Stack.Screen name="(auth)" />
      )} */}
      <Stack.Screen name="(auth)" />
    </Stack>
  );
}
