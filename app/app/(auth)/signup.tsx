import { useState } from "react";
import { Alert, StyleSheet, View } from "react-native";
import { signup } from "@/lib/api";
import { Button, Input } from "@rneui/themed";
import { router } from "expo-router";

export default function SignUp() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSignUp() {
    setLoading(true);
    try {
      await signup(name, email, password);
      Alert.alert("Success", "Check your inbox for verification!");
      router.push("/");
    } catch (error: any) {
      Alert.alert("Error", error.response?.data?.message || "Signup failed");
    }
    setLoading(false);
  }

  return (
    <View style={styles.container}>
      <View style={styles.inputContainer}>
        <Input
          label="Name"
          leftIcon={{ type: "font-awesome", name: "user" }}
          onChangeText={setName}
          value={name}
          placeholder="Joe Hendry"
          autoCapitalize="none"
        />
      </View>
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
          title="Sign Up"
          disabled={loading}
          onPress={handleSignUp}
          buttonStyle={styles.button}
        />
        <Button
          title="Go to Login"
          type="outline"
          onPress={() => router.push("/")}
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
