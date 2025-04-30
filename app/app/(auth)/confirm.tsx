import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from "react-native";
import { Stack, useRouter, useLocalSearchParams } from "expo-router";
import { confirmEmail } from "../../lib/api";
import { StatusBar } from "expo-status-bar";

const ConfirmScreen = () => {
  const router = useRouter();
  const params = useLocalSearchParams();
  const email = params.email as string;

  const [verificationCode, setVerificationCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleConfirm = async () => {
    if (!verificationCode) {
      Alert.alert("Error", "Please enter the verification code");
      return;
    }

    try {
      setIsLoading(true);
      await confirmEmail(email, verificationCode);
      Alert.alert("Success", "Your email has been verified successfully!", [
        { text: "OK", onPress: () => router.replace("/") },
      ]);
    } catch (error: any) {
      Alert.alert(
        "Error",
        error.response?.data?.message || "Failed to verify email"
      );
    } finally {
      setIsLoading(false);
    }
  };

  const resendCode = async () => {
    Alert.alert("Info", "Resend code functionality not implemented yet");
  };

  return (
    <View style={styles.container}>
      <StatusBar style="dark" />
      <Stack.Screen
        options={{
          title: "Verify Email",
          headerStyle: {
            backgroundColor: "#f9fafb",
          },
          headerShadowVisible: false,
        }}
      />

      <View style={styles.formContainer}>
        <Text style={styles.title}>Verify Your Email</Text>
        <Text style={styles.subtitle}>
          We've sent a verification code to{"\n"}
          <Text style={styles.emailText}>{email || "your email"}</Text>
        </Text>

        <View style={styles.inputContainer}>
          <TextInput
            style={styles.input}
            placeholder="Enter verification code"
            value={verificationCode}
            onChangeText={setVerificationCode}
            keyboardType="number-pad"
            maxLength={6}
            autoCapitalize="none"
          />
        </View>

        <TouchableOpacity
          style={styles.button}
          onPress={handleConfirm}
          disabled={isLoading}
        >
          {isLoading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.buttonText}>Verify Email</Text>
          )}
        </TouchableOpacity>

        <View style={styles.footer}>
          <Text style={styles.footerText}>Didn't receive the code? </Text>
          <TouchableOpacity onPress={resendCode}>
            <Text style={styles.linkText}>Resend</Text>
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#f9fafb",
  },
  formContainer: {
    flex: 1,
    padding: 20,
    justifyContent: "center",
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 12,
    color: "#111827",
  },
  subtitle: {
    fontSize: 16,
    color: "#6b7280",
    marginBottom: 36,
    lineHeight: 22,
  },
  emailText: {
    fontWeight: "600",
    color: "#111827",
  },
  inputContainer: {
    marginBottom: 24,
  },
  input: {
    backgroundColor: "#ffffff",
    borderWidth: 1,
    borderColor: "#d1d5db",
    borderRadius: 8,
    padding: 16,
    fontSize: 16,
    color: "#111827",
  },
  button: {
    backgroundColor: "#3b82f6",
    borderRadius: 8,
    padding: 16,
    alignItems: "center",
    marginBottom: 16,
  },
  buttonText: {
    color: "#ffffff",
    fontSize: 16,
    fontWeight: "600",
  },
  footer: {
    flexDirection: "row",
    justifyContent: "center",
    marginTop: 24,
  },
  footerText: {
    color: "#6b7280",
    fontSize: 14,
  },
  linkText: {
    color: "#3b82f6",
    fontSize: 14,
    fontWeight: "600",
  },
});

export default ConfirmScreen;
