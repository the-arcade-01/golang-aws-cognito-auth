import { useState } from "react";
import { Alert, StyleSheet, View } from "react-native";
import { login } from "../../lib/api";
import { Button, Input } from "@rneui/themed";
import { router } from "expo-router";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleLogin() {
    setLoading(true);
    try {
      const data = await login(email, password);
      const token = data.access_token;

      router.replace({
        pathname: "/(tabs)/home",
        params: { newLogin: "true" },
      });
    } catch (error: any) {
      console.log(error);
      Alert.alert("Error", error.response?.data?.message || "Login failed");
    }
    setLoading(false);
  }

  return (
    <View style={styles.container}>
      <View style={styles.inputContainer}>
        <Input
          label="Email"
          leftIcon={{ type: "font-awesome", name: "envelope" }}
          onChangeText={setEmail}
          value={email}
          placeholder="email@address.com"
          autoCapitalize="none"
        />
      </View>
      <View style={styles.inputContainer}>
        <Input
          label="Password"
          leftIcon={{ type: "font-awesome", name: "lock" }}
          onChangeText={setPassword}
          value={password}
          secureTextEntry
          placeholder="Password"
          autoCapitalize="none"
        />
      </View>
      <View style={styles.buttonContainer}>
        <Button
          title="Sign In"
          disabled={loading}
          onPress={handleLogin}
          buttonStyle={styles.button}
        />
        <Button
          title="Go to Sign Up"
          type="outline"
          onPress={() => router.push("/signup")}
          buttonStyle={styles.button}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 12,
    justifyContent: "center",
  },
  inputContainer: {
    marginVertical: 8,
  },
  buttonContainer: {
    marginVertical: 16,
  },
  button: {
    marginVertical: 8,
  },
});
